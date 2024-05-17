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
