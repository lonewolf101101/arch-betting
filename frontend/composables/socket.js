let socket = null

export const useQueue = (callback) => {
  if (socket == null) {
    if (useRuntimeConfig().public.dev) {
      socket = new WebSocket("ws://localhost:4000/api/ws")
    } else {
      socket = new WebSocket("wss://ocr.chimege.com/api/ws")
    }
    socket.onmessage = onSocketMessage
  }

  const queue = useState("queue", () => 0)
  watch(queue, callback)
}

function onSocketMessage(event) {
  const queue = useState("queue", () => 0)
  let data = JSON.parse(event.data)
  if (data.Type === "PING") {
    socket.send(JSON.stringify({ Type: "PONG", Text: ".." }))
    return
  }
  queue.value = data
}

function onSocketSendMessage(data) {
  socket.send(data)
}

export const sendMessage = (data) => {
  return onSocketSendMessage(data)
}
