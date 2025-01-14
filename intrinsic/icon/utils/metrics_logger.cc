// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/utils/metrics_logger.h"

#include <string>

#include "absl/status/status.h"
#include "google/protobuf/struct.pb.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/rt_thread.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::icon {

MetricsLogger::MetricsLogger(std::string module_name)
    : module_name_(module_name) {}

MetricsLogger::~MetricsLogger() {
  if (metrics_publisher_thread_.joinable()) {
    shutdown_requested_.store(true);
    metrics_publisher_thread_.join();
  }
}

absl::Status MetricsLogger::Start() {
  if (metrics_publisher_thread_.joinable()) {
    return absl::FailedPreconditionError(
        "Metrics publisher thread is already running");
  }

  intrinsic::ThreadOptions options;
  options.SetNormalPriorityAndScheduler();
  options.SetName("metrics_publisher_thread_");
  shutdown_requested_.store(false);
  INTR_ASSIGN_OR_RETURN(
      metrics_publisher_thread_,
      CreateRealtimeCapableThread(options, [this]() { LoggerFunction(); }));
  return absl::OkStatus();
}

void MetricsLogger::LoggerFunction() {
  // Read the queue until empty
  while (!shutdown_requested_.load()) {
    return;
  }
}

}  // namespace intrinsic::icon
