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

function headers(auth = false) {
  const h = { "Content-Type": "application/json" }
  if (auth && token) {
    h["Authorization"] = "Bearer " + token
  }
  return h
}

function show(res) {
  document.getElementById("output").innerText =
    JSON.stringify(res, null, 2)
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
  show(await res.json())
}

async function login() {
  const res = await fetch(API + "/login", {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({
      username: loginUsername.value,
      password: loginPassword.value
    })
  })

  const data = await res.json()
  token = data.token
  localStorage.setItem("token", token)
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
  show(await res.json())
}

async function startServer() {
  const res = await fetch(API + "/start", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      server_id: startServerId.value
    })
  })
  show(await res.json())
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
  show(await res.json())
}

async function deleteServer() {
  const res = await fetch(API + "/delete", {
    method: "POST",
    headers: headers(true),
    body: JSON.stringify({
      server_id: deleteServerId.value
    })
  })
  show(await res.text())
}
