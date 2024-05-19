import {DisplayTodoList} from "./DisplayTodoList.jsx";
import {createContext, useState} from "react";
import AddTodo from "./AddTodo.jsx";

export const TodosContext = createContext()

export function TodoList() {
    const [todos, setTodos] = useState([]);
    return (
        <TodosContext.Provider value={setTodos}>
            <DisplayTodoList todos={todos}/>
            <AddTodo/>
        </TodosContext.Provider>
    )
}
