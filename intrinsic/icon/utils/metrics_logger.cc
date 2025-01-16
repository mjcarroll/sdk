// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/utils/metrics_logger.h"

#include <string>

#include "absl/log/log.h"
#include "absl/status/status.h"
#include "google/protobuf/struct.pb.h"
#include "intrinsic/icon/utils/realtime_metrics.h"
#include "intrinsic/logging/data_logger_client.h"
#include "intrinsic/logging/proto/log_item.pb.h"
#include "intrinsic/performance/analysis/proto/performance_metrics.pb.h"
#include "intrinsic/platform/common/buffers/realtime_write_queue.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/util/thread/rt_thread.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"

using intrinsic_proto::data_logger::LogItem;

namespace intrinsic::icon {

MetricsLogger::MetricsLogger(std::string module_name)
    : cycle_time_metrics_queue_(), module_name_(module_name) {}

MetricsLogger::~MetricsLogger() {
  cycle_time_metrics_queue_.Writer().Close();
  shutdown_requested_.store(true);
}

absl::Status MetricsLogger::Start() {
  if (metrics_publisher_thread_.joinable()) {
    return absl::FailedPreconditionError(
        "Metrics publisher thread is already running");
  }

  intrinsic::ThreadOptions options;
  options.SetNormalPriorityAndScheduler();
  options.SetName("metrics_publisher");
  shutdown_requested_.store(false);
  return absl::OkStatus();
}

bool MetricsLogger::AddCycleTimeMetrics(
    const CycleTimeMetrics& cycle_time_metrics) {
  return true;
}
void MetricsLogger::PublishMetrics() {
  // Cyclically read metrics from the queue and log them.
  // Publishes metrics as fast as they are available.
  while (!shutdown_requested_.load()) {
    // No busy loop, because PublishCycleTimeMetrics blocks until metrics are
    // available.
    PublishCycleTimeMetrics();
  }
}

void MetricsLogger::PublishCycleTimeMetrics() {
}

}  // namespace intrinsic::icon
