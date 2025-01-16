// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_
#define INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_

#include <stdbool.h>

#include <atomic>
#include <string>

#include "absl/status/status.h"
#include "intrinsic/icon/testing/realtime_annotations.h"
#include "intrinsic/icon/utils/realtime_metrics.h"
#include "intrinsic/platform/common/buffers/realtime_write_queue.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::icon {

// A non-real time logger that can be used to publish messages from realtime
// contexts.
class MetricsLogger {
 public:
  // Constructs a MetricsLogger that exports metrics to a topic with
  // `module_name`.
  explicit MetricsLogger(std::string module_name);

  ~MetricsLogger();

  // Starts the metrics logger thread
  absl::Status Start();

  // Copies cycle time metrics into the rt to non-rt queue so they can be
  // logged.
  // Returns false if the queue is full.
  bool AddCycleTimeMetrics(const CycleTimeMetrics& cycle_time_metrics)
      INTRINSIC_CHECK_REALTIME_SAFE;

 private:
  // Cyclically called by the non-rt metrics_publisher_thread_ thread to publish
  // metrics.
  void PublishMetrics();

  // Blocks until data is available in the cycle time metrics queue,
  // Converts the data to a PerformanceMetrics proto and logs it to data logger.
  void PublishCycleTimeMetrics();

  // Atomic flag to enable/disable the metrics thread.
  std::atomic<bool> shutdown_requested_;

  RealtimeWriteQueue<CycleTimeMetrics> cycle_time_metrics_queue_;

  //  Thread to publish metrics (non-real-time)
  intrinsic::Thread metrics_publisher_thread_;

  // The name of the module that is logging metrics
  std::string module_name_;
};
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_
