const regUsername = document.getElementById("reg-username")
const regEmail = document.getElementById("reg-email")
const regPassword = document.getElementById("reg-password")

const loginUsername = document.getElementById("login-username")
const loginPassword = document.getElementById("login-password")

const serverVersion = document.getElementById("server-version")

const startServerId = document.getElementById("start-server-id")

const stopServerId = document.getElementById("stop-server-id")
const stopContainerId = document.getElementById("stop-container-id")

const deleteServerId = document.getElementById("delete-server-id")

const API = "/api/v1"
let token = localStorage.getItem("token") || ""
let userID = "" 

function headers(auth = false) {
  const h = { "Content-Type": "application/json" }
  if (auth && token) {
    h["Authorization"] = "Bearer " + token
  }
  return h
}

function show(data) {
  const output = document.getElementById("output")

  if (typeof data === "string") {
    output.innerText = data
  } else if (typeof data === "object" && data !== null) {
    const values = Object.values(data).join(" ")
    output.innerText = values
  } else {
    output.innerText = String(data)
  }
}

function notify(msg, ok = true) {
  alert((ok ? "✅ " : "❌ ") + msg)
}

async function handleResponse(res) {
  const text = await res.text()
  try {
    return JSON.parse(text)
  } catch {
    return { message: text }
  }
}

async function register() {
  const res = await fetch(API + "/register", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({
      username: regUsername.value,
      email: regEmail.value,
      password: regPassword.value
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "register failed", false)
    return
  }

  userID = data.id

  notify("successfully registered. Your ID is: " + userID)
  show(data)
}


async function login() {
  if (!userID) {
    notify("You must register first to get your user ID", false)
    return
  }

  const res = await fetch(API + "/login", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({
      username: loginUsername.value,
      password: loginPassword.value,
      id: userID
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "login failed", false)
    return
  }

  token = data.token
  localStorage.setItem("token", token)

  notify("login successful")
  show(data)
}

async function createServer() {
  const res = await fetch(API + "/create", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      version: serverVersion.value
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "create failed", false)
    return
  }

  notify("server created")
  show(data)
}

async function startServer() {
  const res = await fetch(API + "/start", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      server_id: startServerId.value
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "start failed", false)
    return
  }

  notify("server started")
  show(data)
}

async function stopServer() {
  const res = await fetch(API + "/stop", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      server_id: stopServerId.value,
      container_id: stopContainerId.value
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "stop failed", false)
    return
  }

  notify("server stopped")
}

async function deleteServer() {
  const res = await fetch(API + "/delete", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      server_id: deleteServerId.value
    })
  })

  const data = await handleResponse(res)

  if (!res.ok) {
    notify(data.message || "delete failed", false)
    return
  }

  notify("server deleted")
}

function copyOutput() {
    const output = document.getElementById("output").innerText
    if (!output) return alert("Nothing to copy!")
    navigator.clipboard.writeText(output)
        .then(() => alert("✅ Copied to clipboard!"))
        .catch(err => alert("❌ Failed to copy: " + err))
}