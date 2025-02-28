# Copyright 2023 Intrinsic Innovation LLC

# Home of pubsub, a non-real-time Publish-Subscribe messaging
# for robotics systems.

load("@bazel_skylib//rules:common_settings.bzl", "string_flag")

package(default_visibility = ["//visibility:public"])

string_flag(
    name = "impl",
    build_setting_default = "zenoh",
    values = [
        "zenoh",
    ],
)

config_setting(
    name = "zenoh_build",
    flag_values = {":impl": "zenoh"},
)

cc_library(
    name = "zenoh_pubsub_data",
    srcs = ["zenoh_pubsub_data.cc"],
    hdrs = ["zenoh_pubsub_data.h"],
    deps = [
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_config",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/base:no_destructor",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/synchronization",
    ],
)

# PubSub API library. The user-facing API providing access to
# publish-subscribe system.
cc_library(
    name = "pubsub",
    srcs = select({
        ":zenoh_build": [
            "zenoh_pubsub.cc",
        ],
    }),
    hdrs = [
        "pubsub.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":kvstore",
        ":publisher",
        ":queryable",
        ":subscription",
        ":zenoh_publisher_data",
        ":zenoh_pubsub_data",
        ":zenoh_subscription_data",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_config",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/base",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@com_google_protobuf//:protobuf",
    ],
)

# PubSub API library. The user-facing API providing access to
# publish-subscribe system.
cc_library(
    name = "kvstore",
    srcs = ["kvstore.cc"],
    hdrs = ["kvstore.h"],
    deps = [
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_helpers",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/strings:string_view",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "zenoh_publisher_data",
    hdrs = ["zenoh_publisher_data.h"],
)

cc_library(
    name = "publisher",
    # We need a different publisher implementation also when we move to a
    # different pubsub implementation
    srcs = select({
        ":zenoh_build": ["zenoh_publisher.cc"],
    }),
    hdrs = [
        "publisher.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":publisher_stats",
        ":zenoh_publisher_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/strings:string_view",
        "@com_google_absl//absl/time",
        "@com_google_protobuf//:protobuf",
        "@io_opencensus_cpp//opencensus/stats",
    ],
)

cc_library(
    name = "queryable",
    # We need a different publisher implementation also when we move to a
    # different pubsub implementation
    srcs = select({
        ":zenoh_build": ["zenoh_queryable.cc"],
    }),
    hdrs = [
        "queryable.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "zenoh_subscription_data",
    hdrs = ["zenoh_subscription_data.h"],
    deps = [
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
    ],
)

cc_library(
    name = "subscription",
    # We need a different subscription implementation also when we move to a
    # different pubsub implementation
    srcs = select({
        ":zenoh_build": [
            "zenoh_subscription.cc",
        ],
    }),
    hdrs = [
        "subscription.h",
    ],
    copts = [
        "-fexceptions",
    ],
    features = ["-use_header_modules"],
    deps = [
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "publisher_stats",
    srcs = ["publisher_stats.cc"],
    hdrs = ["publisher_stats.h"],
    deps = [
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/debugging:leak_check",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
    ],
)

cc_library(
    name = "zenoh_pubsub",
    srcs = ["zenoh_pubsub.cc"],
    deps = [
        ":kvstore",
        ":publisher",
        ":pubsub",
        ":queryable",
        ":subscription",
        ":zenoh_publisher_data",
        ":zenoh_pubsub_data",
        ":zenoh_subscription_data",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_config",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "zenoh_queryable",
    srcs = ["zenoh_queryable.cc"],
    deps = [
        ":queryable",
        "//intrinsic/platform/pubsub/adapters:pubsub_cc_proto",
        "//intrinsic/platform/pubsub/zenoh_util:zenoh_handle",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings:str_format",
    ],
)
