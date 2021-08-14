#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# File: diff_test_result.py
# Created Date: 2021-08-14
# Author: Rabbit
# --------------------------------
# Copyright (c) 2021 Rabbit

import os
import sys
import difflib

ROOT_DIR = os.path.join(os.path.dirname(__file__), "../")
TEST_DIR = os.path.join(ROOT_DIR, "test/")
LOG_DIR = os.path.join(ROOT_DIR, "log/")

def main():
    if len(sys.argv) < 2:
        print("Usage: %proc [case]")
        sys.exit(1)

    case = sys.argv[1]

    out_filepath = os.path.join(TEST_DIR, '%s.out.txt' % case)
    stdout_filepath = os.path.join(TEST_DIR, '%s.stdout.txt' % case)

    with open(out_filepath, 'r', encoding='utf-8') as fp:
        out_lines = fp.readlines()

    with open(stdout_filepath, 'r', encoding='utf-8') as fp:
        stdout_lines = fp.readlines()

    d = difflib.HtmlDiff()
    html = d.make_file(out_lines, stdout_lines, context=True, numlines=20)

    html_filepath = os.path.join(LOG_DIR, '%s.diff.html' % case)

    with open(html_filepath, 'w', encoding='utf-8') as fp:
        fp.write(html)

if __name__ == "__main__":
    main()
