import {DisplayTodoList} from "./DisplayTodoList.jsx";
import {AddTodo} from "./AddTodo.jsx";
import {createContext, useState} from "react";

export const TodosContext = createContext(undefined)

export function TodoList() {
    const [todos, setTodos] = useState([]);
    return (
        <TodosContext.Provider value={{todos, setTodos}}>
            <DisplayTodoList/>
            <AddTodo/>
        </TodosContext.Provider>
    )
}
