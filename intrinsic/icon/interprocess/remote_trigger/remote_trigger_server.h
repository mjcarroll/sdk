// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_INTERPROCESS_REMOTE_TRIGGER_REMOTE_TRIGGER_SERVER_H_
#define INTRINSIC_ICON_INTERPROCESS_REMOTE_TRIGGER_REMOTE_TRIGGER_SERVER_H_

#include <atomic>
#include <functional>
#include <string>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "intrinsic/icon/interprocess/binary_futex.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/memory_segment.h"
#include "intrinsic/icon/interprocess/shared_memory_manager/shared_memory_manager.h"
#include "intrinsic/util/thread/thread.h"
#include "intrinsic/util/thread/thread_options.h"

namespace intrinsic::icon {

using RemoteTriggerServerCallback = std::function<void(void)>;

// A RemoteTriggerServer listens to incoming requests from a client and executes
// its callback when a request is issued.
// Remote denotes support for inter-process communication, yet still
// requires the connection to be machine-local, meaning the two processes for
// the client and server have to be executed on the same computer. The
// connection between a server and client is based on a named semaphore. The
// `server_id` passed into the server and client thus have to match in order to
// establish a connection.
// There is a recommended 1:1 relationship between a server and a client;
// the server can't distinguish a request when being triggered by various
// clients. While we can make sure that one client can only trigger one request
// at the time, we can't easily prevent multiple clients (each in their own
// process) to trigger a request at the same time.
class RemoteTriggerServer final {
 public:
  // Creates a new server instance named `server_memory_name` on `shm_manager`.
  // When the server is signaled, it executes the callback and signals a
  // response back to the client when done.
  static absl::StatusOr<RemoteTriggerServer> Create(
      intrinsic::icon::SharedMemoryManager& shm_manager,
      absl::string_view server_memory_name,
      RemoteTriggerServerCallback&& callback);

  // This class is move-only.
  RemoteTriggerServer(const RemoteTriggerServer& other) = delete;
  RemoteTriggerServer& operator=(const RemoteTriggerServer& other) = delete;
  // Moving an instance will stop the server. If the server was previously
  // running, it has to be explicitly restarted afterwards.
  RemoteTriggerServer(RemoteTriggerServer&& other) noexcept;
  RemoteTriggerServer& operator=(RemoteTriggerServer&& other) noexcept;
  // Destructor.
  ~RemoteTriggerServer();

  // Starts the server loop within the current thread.
  // This call blocks indefinitely until `Stop()` is called from another
  // thread. If the server was previously already started, another call to
  // `Start()` returns immediately.
  // Given its blocking behavior, this function should be used with
  // external thread handling.
  void Start();

  // Starts the server loop within a new thread.
  // Unlike `Start()`, the function returns immediately, running the server loop
  // in an internal thread with specified thread options.
  // The thread can then be stopped with a call to `Stop()`. A call to
  // `StartAsync()` has no effect if the server is already started.
  // Returns `absl::OkStatus()` upon success, error status when thread could not
  // start correctly.
  absl::Status StartAsync(const intrinsic::ThreadOptions& thread_options = {});

  // Queries whether the server has started.
  bool IsStarted() const;

  // Stops the current server loop.
  // A call to `RequestStop()` exits the server loop independently whether it's
  // been started via `Start()` or `StartAsync()`. A call to `RequestStop()` has
  // no effect if the server is already stopped.
  //
  // Use this to enter a "lame duck" mode and ensure the thread doesn't start
  // any further requests. This can be useful on shutdown, where a request might
  // get stuck until you call the corresponding HWM's `Shutdown()` method.
  //
  // Remember to call `JoinAsyncThread()` if you want to restart the server! If
  // not, you can rely on the destructor to join the async thread, if any.
  void RequestStop();

  // If there is an async server thread (i.e. the server was started with
  // `StartAsync()`), this joins that thread.
  // This can block if there's an ongoing call, so be sure that there isn't, or
  // that you have a way to unblock ongoing calls (for hardware modules, the
  // `Shutdown()` method *should* do that).
  //
  // If there is no async server thread, or it's already joined, this is a
  // no-op.
  void JoinAsyncThread();

  // Queries whether the server is ready to start.
  // This returns true if the server is stopped and any asynchronous threads
  // have already been joined.
  bool IsReadyToStart() const;

  // Queries the server once and executes the callback if a request is ready.
  // Does not execute the callback if the server is started already.
  // Returns true if a callback was triggered, false if not.
  bool Query();

 private:
  // Main loop function.
  // Waits for an incoming trigger sent by a client and calls the provided
  // callback upon arrival. Once the callback returns, it sends a response
  // notification to the client, indicating that the callback has been
  // completed.
  void Run();

  RemoteTriggerServer(absl::string_view server_memory_name,
                      RemoteTriggerServerCallback&& callback,
                      ReadOnlyMemorySegment<BinaryFutex>&& request_futex,
                      ReadWriteMemorySegment<BinaryFutex>&& response_futex);

  std::string server_memory_name_;
  RemoteTriggerServerCallback callback_;
  // initialize to `false`, indicating the system is currently stopped.
  std::atomic<bool> is_running_{false};
  // The interprocess signaling is done via two semaphores shared between a
  // server and its clients.
  ReadOnlyMemorySegment<BinaryFutex> request_futex_;
  ReadWriteMemorySegment<BinaryFutex> response_futex_;

  intrinsic::Thread async_thread_;
};

}  // namespace intrinsic::icon

#endif  // INTRINSIC_ICON_INTERPROCESS_REMOTE_TRIGGER_REMOTE_TRIGGER_SERVER_H_
