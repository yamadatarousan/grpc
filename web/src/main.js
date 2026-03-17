import { createClient } from "@connectrpc/connect";
import { createGrpcWebTransport } from "@connectrpc/connect-web";
import { TodoService } from "./gen/todo/v1/todo_pb";
import "./style.css";
const transport = createGrpcWebTransport({
    baseUrl: window.location.origin,
});
const client = createClient(TodoService, transport);
const form = document.getElementById("add-form");
const input = document.getElementById("title-input");
const list = document.getElementById("todo-list");
const reloadButton = document.getElementById("reload-button");
const status = document.getElementById("status");
function renderStatus(message) {
    status.textContent = message;
}
async function reloadTodos() {
    try {
        renderStatus("");
        const response = await client.listTodos({});
        list.innerHTML = "";
        for (const todo of response.todos) {
            const li = document.createElement("li");
            const createdAt = new Date(Number(todo.createdAtUnix) * 1000).toLocaleString();
            li.textContent = `#${todo.id} ${todo.title} (${createdAt})`;
            list.appendChild(li);
        }
        if (response.todos.length === 0) {
            const li = document.createElement("li");
            li.textContent = "TODOはまだありません";
            list.appendChild(li);
        }
    }
    catch (err) {
        renderStatus(`一覧取得に失敗: ${String(err)}`);
    }
}
form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const title = input.value.trim();
    if (!title) {
        renderStatus("タイトルを入力してください");
        return;
    }
    try {
        renderStatus("");
        await client.addTodo({ title });
        input.value = "";
        await reloadTodos();
    }
    catch (err) {
        renderStatus(`追加に失敗: ${String(err)}`);
    }
});
reloadButton.addEventListener("click", async () => {
    await reloadTodos();
});
void reloadTodos();
