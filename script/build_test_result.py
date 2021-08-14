#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# File: build_test_result.py
# Created Date: 2021-08-14
# Author: Rabbit
# --------------------------------
# Copyright (c) 2021 Rabbit

import os
import re
import jieba

ROOT_DIR = os.path.join(os.path.dirname(__file__), "../")
TEST_DIR = os.path.join(ROOT_DIR, "test/")

re_input = re.compile(r'^(.*)\.in\.txt$')

def dump_result(filepath, res):
    with open(filepath, 'w', encoding='utf-8') as fp:
        for word in res:
            fp.write(word)
            fp.write('\n')

def main():
    for filename in os.listdir(TEST_DIR):
        mat = re_input.match(filename)
        if mat is None:
            continue
        in_filepath = os.path.join(TEST_DIR, filename)

        with open(in_filepath, 'r', encoding='utf-8') as fp:
            content = fp.read()

        out_filepath = os.path.join(TEST_DIR, '%s.accurate.stdout.txt' % mat.group(1))
        dump_result(out_filepath, jieba.cut(content))

        out_filepath = os.path.join(TEST_DIR, '%s.all.stdout.txt' % mat.group(1))
        dump_result(out_filepath, jieba.cut(content, cut_all=True))

        out_filepath = os.path.join(TEST_DIR, '%s.search.stdout.txt' % mat.group(1))
        dump_result(out_filepath, jieba.cut_for_search(content))

if __name__ == "__main__":
    main()
