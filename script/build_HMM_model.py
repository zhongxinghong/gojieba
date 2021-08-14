#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# File: build_finalseg_model.py
# Created Date: 2021-08-14
# Author: Rabbit
# --------------------------------
# Copyright (c) 2021 Rabbit

import os
import sys

JIEBA_ROOT_DIR = os.path.join(ROOT_DIR, "jieba/")
ROOT_DIR = os.path.join(os.path.dirname(__file__), "../")

PYJIEBA_FINALSEG_MODEL_DIR = os.path.join(JIEBA_ROOT_DIR, "jieba/finalseg/")
GOJIEBA_FINALSEG_MODEL_DIR = os.path.join(ROOT_DIR, "internal/hmm/model/")

def build_gojieba_finalseg_model():
    sys.path.append(PYJIEBA_FINALSEG_MODEL_DIR)

    import prob_start
    import prob_trans
    import prob_emit

    prob_start_dst = os.path.join(GOJIEBA_FINALSEG_MODEL_DIR, "prob_start.txt")
    prob_trans_dst = os.path.join(GOJIEBA_FINALSEG_MODEL_DIR, "prob_trans.txt")
    prob_emit_dst = os.path.join(GOJIEBA_FINALSEG_MODEL_DIR, "prob_emit.txt")

    with open(prob_start_dst, 'w', encoding='utf-8') as fp:
        for state, prob in prob_start.P.items():
            assert len(state) == 1
            fp.write("%s %s\n" % (state, prob))

    with open(prob_trans_dst, 'w', encoding='utf-8') as fp:
        for prev, d in prob_trans.P.items():
            assert len(prev) == 1
            for state, prob in d.items():
                assert len(state) == 1
                fp.write("%s %s %s\n" % (prev, state, prob))

    with open(prob_emit_dst, 'w', encoding='utf-8') as fp:
        for state, d in prob_emit.P.items():
            assert len(state) == 1
            for word, prob in d.items():
                assert len(word) == 1
                fp.write("%s %s %s\n" % (state, word, prob))

def main():
    build_gojieba_finalseg_model()

if __name__ == "__main__":
    main()
