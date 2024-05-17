import {createContext, useContext, useEffect, useState} from "react";
import {apiUrl} from "../App.jsx";
import {TodosContext} from "./TodoList.jsx";


export function DisplayTodoList() {
    const {todos, setTodos} = useContext(TodosContext)

    useEffect(() => {
        getTodos()
        const interval = setInterval(getTodos, 5000)
        return () => {
            clearInterval(interval)
        }
    }, [])

    const getTodos = () => {
        fetch(apiUrl).then((res) => {
            res.json().then((data) => {
                setTodos(data);
            })
        })
    }
    const updateTodo = async (todo) => {
        fetch(`${apiUrl}/${todo.id}`,
            {method: "PATCH", body: JSON.stringify(todo)})
    }
    const deleteTodo = async (id) => {
        fetch(`${apiUrl}/${id}`, {method: "DELETE"})
    }

    const checkHandler = (id) => {
        setTodos((prevTodos) => prevTodos.map((todo) => {
                if (todo.id === id) {
                    let newTodo = {...todo, done: !todo.done};
                    updateTodo(newTodo)
                    return newTodo;
                }
                return todo;
            }
        ));
    }
    const deleteHandler = (id) => {
        setTodos((prevTodos) => prevTodos.filter((todo) => todo.id !== id))
        deleteTodo(id)
    }

    return (
        <>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 m-4">
                {todos.map((todo) => (
                    <div key={todo.id} className="bg-gray-300  rounded-lg  p-4 border-2 border-black ">
                        <h2 className="text-xl font-semibold ">{todo.title}</h2>
                        <div className="mt-2 flex items-center">
                            <input type="checkbox" checked={todo.done}
                                   className=" h-4 w-4 mr-2" onClick={() => {
                                checkHandler(todo.id);
                            }} readOnly/>
                            <span className="text-sm">{todo.done ? 'Done' : 'Not done'}</span>
                        </div>
                        <button onClick={() => deleteHandler(todo.id)}
                                className="bg-red-600 rounded p-2 text-white font-semibold mt-4">
                            Delete todo
                        </button>
                    </div>
                ))}
            </div>
        </>
    )
}
