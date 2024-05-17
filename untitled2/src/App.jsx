import {DisplayTodoList} from "./components/DisplayTodoList.jsx";
import {AddTodo} from "./components/AddTodo.jsx";
import {createContext} from "react";
import {TodoList} from "./components/TodoList.jsx";

export const apiUrl = "http://localhost:8080/todos"

function App() {
    return (
        <>
            <TodoList/>
        </>
    )
}

export default App
