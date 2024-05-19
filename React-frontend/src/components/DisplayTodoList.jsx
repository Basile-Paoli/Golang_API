import {useContext, useEffect} from "react";
import {apiUrl} from "../App.jsx";
import {TodosContext} from "./TodoList.jsx";
import PropTypes from "prop-types";


export function DisplayTodoList({todos}) {
    const setTodos = useContext(TodosContext)

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
                setTodos((currentTodos) => {
                    currentTodos.sort((a, b) => a.id - b.id)
                    data.sort((a, b) => a.id - b.id)
                    if (JSON.stringify(currentTodos) !== JSON.stringify(data)) {
                        return data
                    }
                    return currentTodos
                })
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
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4 m-4 min-w-72">
                {todos.map((todo) => (
                    <div key={todo.id} className="bg-gray-300  rounded-lg  p-4 border-2 border-black min-w-32">
                        <h2 className="text-xl font-semibold ">{todo.title}</h2>
                        <div className="mt-2 flex items-center">
                            <input type="checkbox" checked={todo.done}
                                   className=" h-4 w-4 mr-2" onClick={() => {
                                checkHandler(todo.id);
                            }} readOnly/>
                            <span className="text-sm">{todo.done ? 'Done' : 'Pending'}</span>
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

DisplayTodoList.propTypes = {
    todos: PropTypes.array.isRequired
}