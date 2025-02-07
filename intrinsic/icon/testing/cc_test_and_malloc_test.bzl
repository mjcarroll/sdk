# Copyright 2023 Intrinsic Innovation LLC

"""Malloc test is not implemented externally, only invoke regular cc_test."""

load("@bazel_skylib//lib:new_sets.bzl", "sets")

def cc_test_and_malloc_test(name, deps = [], local_defines = [], tags = [], **kwargs):
    native.cc_test(
        name = name,
        local_defines = local_defines,
        tags = tags,
        deps = deps + [
            "//intrinsic/util/testing:gtest_wrapper_main",
        ],
        **kwargs
    )
