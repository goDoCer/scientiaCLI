#!/usr/bin/env python3
"""
scientia-pull — download a student's own coursework submissions, specs,
marks, and results across any academic year(s) from Scientia
(scientia.doc.ic.ac.uk).

The upstream scientiaCLI (Go) only downloads course materials for the
current year. This script targets the /me/ and /{year}/students/ endpoints
that expose the user's own uploaded submissions and grades, and supports
arbitrary years.

Usage:
    python3 pull_submissions.py login
    python3 pull_submissions.py pull --years 2122 2223 2425 2526 --out ./archive
    python3 pull_submissions.py pull --years 2526             # current year only

Endpoints used (from scientia.doc.ic.ac.uk/api/openapi.json):
    POST /api/auth/login
    POST /api/auth/refresh
    GET  /api/me/{year}/exercises
    GET  /api/me/{year}/{module_code}/exercises/{exercise_number}
    GET  /api/{year}/{module_code}/exercises/{exercise_number}/spec
    GET  /api/{year}/{module_code}/exercises/{exercise_number}/submissions/{submission_id}/file
    GET  /api/{year}/{module_code}/exercises/{exercise_number}/model-answer
    GET  /api/{year}/{module_code}/exercises/{exercise_number}/supplementary
    GET  /api/{year}/{username}/exercise-summaries
    GET  /api/{year}/students/{username}/results
"""

import argparse
import getpass
import json
import os
import re
import sys
import threading
import urllib.error
import urllib.parse
import urllib.request
from concurrent.futures import ThreadPoolExecutor, as_completed
from http.cookiejar import CookieJar
from pathlib import Path

BASE = "https://scientia.doc.ic.ac.uk/api"
CONFIG_DIR = Path.home() / ".scientia-cli"
TOKEN_FILE = CONFIG_DIR / "tokens.json"


def load_tokens():
    if not TOKEN_FILE.exists():
        return {}
    return json.loads(TOKEN_FILE.read_text())


def save_tokens(tokens):
    CONFIG_DIR.mkdir(exist_ok=True)
    TOKEN_FILE.write_text(json.dumps(tokens, indent=2))
    TOKEN_FILE.chmod(0o600)


def build_cookie_header(tokens):
    parts = []
    if tokens.get("access_token"):
        parts.append(f"access_token_cookie={tokens['access_token']}")
    if tokens.get("refresh_token"):
        parts.append(f"refresh_token_cookie={tokens['refresh_token']}")
    return "; ".join(parts)


def api_request(method, path, tokens, body=None, accept_json=True, retry=True):
    """Perform an authenticated request. Returns (status, headers, body_bytes)."""
    url = f"{BASE}{path}"
    data = None
    headers = {"Cookie": build_cookie_header(tokens)}
    if accept_json:
        headers["Accept"] = "application/json"
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            return resp.status, resp.headers, resp.read()
    except urllib.error.HTTPError as e:
        if e.code == 401 and retry and tokens.get("refresh_token"):
            if refresh(tokens):
                return api_request(method, path, tokens, body, accept_json, retry=False)
        return e.code, e.headers, e.read()


def refresh(tokens):
    """Try to refresh access token using refresh_token_cookie."""
    req = urllib.request.Request(
        f"{BASE}/auth/refresh",
        method="POST",
        headers={"Cookie": f"refresh_token_cookie={tokens.get('refresh_token','')}"},
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            new = parse_set_cookie(resp.headers.get_all("Set-Cookie") or [])
            if "access_token" in new:
                tokens["access_token"] = new["access_token"]
                if "refresh_token" in new:
                    tokens["refresh_token"] = new["refresh_token"]
                save_tokens(tokens)
                return True
    except urllib.error.HTTPError:
        pass
    return False


def parse_set_cookie(cookies):
    """Extract access_token_cookie and refresh_token_cookie values."""
    out = {}
    for c in cookies:
        m = re.match(r"(access_token|refresh_token)_cookie=([^;]+)", c)
        if m:
            out[m.group(1)] = m.group(2)
    return out


def cmd_login(args):
    shortcode = args.shortcode or os.environ.get("SCIENTIA_SHORTCODE")
    if not shortcode:
        if not sys.stdin.isatty():
            print("No TTY and no --shortcode / SCIENTIA_SHORTCODE env var provided.",
                  file=sys.stderr)
            sys.exit(2)
        shortcode = input("Imperial shortcode: ").strip()

    password = os.environ.get("SCIENTIA_PASSWORD")
    if not password:
        if not sys.stdin.isatty():
            print("No TTY and no SCIENTIA_PASSWORD env var provided.",
                  file=sys.stderr)
            sys.exit(2)
        password = getpass.getpass("Password (hidden): ")

    req = urllib.request.Request(
        f"{BASE}/auth/login",
        method="POST",
        data=json.dumps({"username": shortcode, "password": password}).encode(),
        headers={"Content-Type": "application/json", "Accept": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            cookies = parse_set_cookie(resp.headers.get_all("Set-Cookie") or [])
    except urllib.error.HTTPError as e:
        print(f"Login failed: {e.code} {e.reason}", file=sys.stderr)
        print(e.read().decode(errors="replace"), file=sys.stderr)
        sys.exit(1)

    if "access_token" not in cookies:
        print("Login response did not include access_token cookie", file=sys.stderr)
        sys.exit(1)

    tokens = {
        "shortcode": shortcode,
        "access_token": cookies["access_token"],
        "refresh_token": cookies.get("refresh_token", ""),
    }
    save_tokens(tokens)
    print(f"Logged in as {shortcode}. Tokens saved to {TOKEN_FILE}")


def sanitize(name):
    """Filesystem-safe name."""
    name = re.sub(r'[/\\:*?"<>|]', "-", str(name))
    return name.strip().strip(".") or "unnamed"


def year_label(y):
    """2122 → 2021-22, 2526 → 2025-26."""
    if len(y) == 4 and y.isdigit():
        return f"20{y[:2]}-{y[2:]}"
    return y


def save_json(path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2))


def save_bytes(path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_bytes(data)


def get_json(path, tokens):
    status, _, body = api_request("GET", path, tokens)
    if status != 200:
        return None, (status, body[:200].decode(errors="replace"))
    try:
        return json.loads(body), None
    except json.JSONDecodeError as e:
        return None, ("parse-error", str(e))


def get_bytes(path, tokens):
    status, headers, body = api_request("GET", path, tokens, accept_json=False)
    if status != 200:
        return None, headers, (status, body[:200].decode(errors="replace"))
    return body, headers, None


def cmd_pull(args):
    tokens = load_tokens()
    if not tokens.get("access_token"):
        print("Not logged in. Run: pull_submissions.py login", file=sys.stderr)
        sys.exit(1)
    shortcode = tokens["shortcode"]
    out = Path(args.out).expanduser().resolve()
    out.mkdir(parents=True, exist_ok=True)

    for year in args.years:
        year_dir = out / f"{year_label(year)}"
        year_dir.mkdir(exist_ok=True)
        print(f"\n=== YEAR {year} ({year_label(year)}) ===")

        # Personal results (grades)
        results, err = get_json(f"/{year}/students/{shortcode}/results", tokens)
        if results is not None:
            save_json(year_dir / "results.json", results)
            print(f"  results.json saved")
        else:
            print(f"  results: {err}")

        # Exercise summaries (marks per exercise)
        summaries, err = get_json(f"/{year}/{shortcode}/exercise-summaries", tokens)
        if summaries is not None:
            save_json(year_dir / "exercise-summaries.json", summaries)
            print(f"  exercise-summaries.json saved")
        else:
            print(f"  exercise-summaries: {err}")

        # List all exercises for this student in this year
        exercises, err = get_json(f"/me/{year}/exercises", tokens)
        if exercises is None:
            print(f"  /me/{year}/exercises: {err}")
            continue

        if not isinstance(exercises, list):
            print(f"  /me/{year}/exercises: unexpected shape: {type(exercises)}")
            save_json(year_dir / "exercises-raw.json", exercises)
            continue

        save_json(year_dir / "exercises-list.json", exercises)

        # Filter: only exercises the user actually participated in
        # (has submissions OR was marked OR has feedback).
        # /me/{year}/exercises returns everything accessible incl. reference material
        # from prior years — that's not what we want to archive.
        def participated(ex):
            if ex.get("submissions"):
                return True
            if ex.get("mark") is not None:
                return True
            if ex.get("feedback"):
                return True
            return False

        mine = [ex for ex in exercises if participated(ex)]
        print(f"  {len(exercises)} exercises listed, {len(mine)} participated in"
              + (f" (limit={args.limit})" if args.limit else ""))

        picked = mine[: args.limit] if args.limit else mine

        # Prep: figure out dir + normalized module/number/title for each exercise
        prepped = []
        for ex in picked:
            mc = ex.get("module_code") or ex.get("moduleCode") or ex.get("module")
            num = ex.get("number") or ex.get("exercise_number") or ex.get("exerciseNumber")
            title = ex.get("title") or ex.get("name") or f"ex-{num}"
            if not mc or num is None:
                print(f"    ! skip (missing module/number): {ex}")
                continue
            mod_dir = year_dir / sanitize(mc)
            ex_dir = (mod_dir / f"{num:02d}-{sanitize(title)}"
                      if isinstance(num, int)
                      else mod_dir / f"{sanitize(str(num))}-{sanitize(title)}")
            ex_dir.mkdir(parents=True, exist_ok=True)
            prepped.append((mc, num, title, ex_dir))

        # Phase A: fetch details.json for every exercise in parallel
        details_map = {}
        with ThreadPoolExecutor(max_workers=args.workers) as pool:
            futures = {
                pool.submit(get_json, f"/me/{year}/{mc}/exercises/{num}", tokens): (mc, num, ex_dir)
                for (mc, num, _, ex_dir) in prepped
            }
            for fut in as_completed(futures):
                mc, num, ex_dir = futures[fut]
                details, err = fut.result()
                details_map[(mc, num)] = (details, err)
                if details is not None:
                    save_json(ex_dir / "details.json", details)

        # Phase B: build all file-fetch tasks from details, then run them in parallel
        tasks = []  # list of (mc, num, ex_dir, kind, path, out_name_hint_or_None)
        for (mc, num, title, ex_dir) in prepped:
            details, _ = details_map.get((mc, num), (None, None))

            if details and details.get("spec"):
                tasks.append((mc, num, ex_dir, "spec",
                              f"/{year}/{mc}/exercises/{num}/spec", None))

            for s in (details or {}).get("submissions") or []:
                sid = s.get("id")
                fname = s.get("target_submission_file_name") or f"submission-{sid}.bin"
                if sid is None:
                    continue
                tasks.append((mc, num, ex_dir, f"sub#{sid}",
                              f"/{year}/{mc}/exercises/{num}/submissions/{sid}/file",
                              sanitize(fname)))

            if details and details.get("model_answer"):
                tasks.append((mc, num, ex_dir, "model-answer",
                              f"/{year}/{mc}/exercises/{num}/model-answer", None))

            if details and details.get("supplementary_file"):
                tasks.append((mc, num, ex_dir, "supplementary",
                              f"/{year}/{mc}/exercises/{num}/supplementary", None))

        results_per_ex = {(mc, num): {"got": [], "fails": []} for (mc, num, _, _) in prepped}
        counter_lock = threading.Lock()
        done_count = [0]
        total = len(tasks)

        def detect_ext(body, ct):
            # magic-byte sniff first (more reliable than Content-Type)
            if body[:4] == b"%PDF":
                return ".pdf"
            if body[:4] == b"PK\x03\x04":
                return ".zip"
            if body[:2] == b"\xff\xd8":
                return ".jpg"
            if body[:8] == b"\x89PNG\r\n\x1a\n":
                return ".png"
            ct = (ct or "").lower()
            if "pdf" in ct:
                return ".pdf"
            if "zip" in ct or "octet-stream" in ct:
                return ".zip"
            if "html" in ct:
                return ".html"
            if "presentation" in ct or "powerpoint" in ct:
                return ".pptx"
            if "sheet" in ct or "excel" in ct:
                return ".xlsx"
            return ".bin"

        def do_task(t):
            mc, num, ex_dir, kind, path, out_name = t
            body, hdrs, err = get_bytes(path, tokens)
            if body is None:
                return (mc, num, kind, None, err)

            if out_name:
                save_bytes(ex_dir / out_name, body)
                label = out_name
            else:
                ext = detect_ext(body, hdrs.get("Content-Type"))
                fname = f"{kind}{ext}"
                save_bytes(ex_dir / fname, body)
                label = fname
            return (mc, num, kind, label, None)

        print(f"  fetching {total} files with {args.workers} workers...")
        with ThreadPoolExecutor(max_workers=args.workers) as pool:
            futures = [pool.submit(do_task, t) for t in tasks]
            for fut in as_completed(futures):
                mc, num, kind, label, err = fut.result()
                if err:
                    results_per_ex[(mc, num)]["fails"].append(f"{kind}={err[0]}")
                else:
                    results_per_ex[(mc, num)]["got"].append(label)
                with counter_lock:
                    done_count[0] += 1
                    if done_count[0] % 10 == 0 or done_count[0] == total:
                        print(f"    [{done_count[0]}/{total}]")

        # Summary per exercise
        for (mc, num, title, _) in prepped:
            details, _ = details_map.get((mc, num), (None, None))
            mark = (details or {}).get("mark")
            mark_str = f" mark={mark}" if mark is not None else ""
            r = results_per_ex[(mc, num)]
            got_str = ",".join(["details"] + r["got"]) if details is not None else ",".join(r["got"])
            print(f"    {mc} #{num} {title!r}{mark_str}: got=[{got_str}]"
                  + (f" fails=[{','.join(r['fails'])}]" if r["fails"] else ""))

    print(f"\nDone. Archive at: {out}")


def main():
    p = argparse.ArgumentParser(prog="scientia-pull",
                                 description="Download Imperial DoC coursework submissions/marks/results.")
    sub = p.add_subparsers(dest="cmd", required=True)

    sp_login = sub.add_parser("login", help="log in and save tokens")
    sp_login.add_argument("--shortcode", help="Imperial shortcode (optional; will prompt if omitted)")
    sp_login.set_defaults(func=cmd_login)

    sp_pull = sub.add_parser("pull", help="download submissions/results for given years")
    sp_pull.add_argument("--years", nargs="+", required=True,
                          help="academic years like 2122 2223 2425 2526")
    sp_pull.add_argument("--out", default="./scientia-archive",
                          help="output directory (default: ./scientia-archive)")
    sp_pull.add_argument("--limit", type=int, default=0,
                          help="only pull first N exercises per year (0 = all). Useful for smoke-test.")
    sp_pull.add_argument("--workers", type=int, default=10,
                          help="parallel HTTP workers for file fetches (default 10)")
    sp_pull.set_defaults(func=cmd_pull)

    args = p.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
