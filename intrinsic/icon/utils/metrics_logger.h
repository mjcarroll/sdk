// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_
#define INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_

#include <stdbool.h>

#include <atomic>
#include <string>

#include "absl/status/status.h"
#include "intrinsic/util/thread/thread.h"

namespace intrinsic::icon {

// A non-real time logger that can be used to log messages from real time
// contexts.
class MetricsLogger {
 public:
  // Constructs a MetricsLogger that exports metrics to a topic with
  // `module_name`.
  explicit MetricsLogger(std::string module_name);

  ~MetricsLogger();

  // Starts the metrics logger thread
  absl::Status Start();

 private:
  // Thread function
  void LoggerFunction();
  //  Thread to publish metrics (non-real-time)
  intrinsic::Thread metrics_publisher_thread_;
  // Atomic flag to enable/disable the metrics thread.
  std::atomic<bool> shutdown_requested_;
  // The name of the module that is logging metrics
  std::string module_name_;
};
}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_UTILS_METRICS_LOGGER_H_
