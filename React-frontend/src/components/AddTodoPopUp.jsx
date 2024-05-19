import {useContext, useEffect, useRef, useState} from "react";
import {TodosContext} from "./TodoList.jsx";
import {apiUrl} from "../App.jsx";
import PropTypes from "prop-types";

export function AddTodoPopUp({toggleDisplay}) {
    const setTodos = useContext(TodosContext)
    const popUpRef = useRef(null);
    const inputRef = useRef(null);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (popUpRef.current && !popUpRef.current.contains(event.target)) {
                toggleDisplay(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return () => {
            document.removeEventListener('mousedown', handleClickOutside)
        }
    }, [toggleDisplay]);


    const [title, setTitle] = useState("");
    const titleInputHandler = (e) => {
        setTitle(e.target.value);
    }

    const addTodo = async (todo) => {

        fetch(apiUrl, {
            method: "POST",
            body: JSON.stringify(todo),
            redirect: "follow"
        })
            .then((res) => res.json()
            .then((data) => {
                setTodos((prevTodos) => prevTodos.map((todo) => {
                    if (todo.id === -1) {
                        return {...todo, id: data.id}
                    }
                    return todo
                }))
            }))
            .catch((error) => {
                console.log(error)
            })
    }
    const validateHandler = () => {
        const todo = {
            id: -1,
            title: title,
            done: false
        }
        setTodos((prevTodos) => [...prevTodos, todo])
        addTodo(todo)
        toggleDisplay(false)

    }

    return (
        <div className="items-center inset-0 p-14 fixed">
            <div className="bg-black opacity-50 inset-0 fixed"></div>
            <div ref={popUpRef} className="rounded bg-white relative border-black border-2 max-w-4xl p-4  min-w-72 m-auto">
                <h1 className="justify-between flex items-center font-bold text-xl">
                    <span className="m-auto">Ajouter une Todo</span>
                    <button onClick={() => toggleDisplay(false)}
                            className="text-3xl hover:scale-110 hover:bg-gray-300 h-8 w-8 flex justify-center items-center">&times;</button>
                </h1>
                <div className=" flex m-auto mt-4 justify-center">
                    <span className="mr-3">Title : </span>
                    <input onChange={titleInputHandler} value={title} placeholder="Title" autoFocus ref={inputRef}
                           onKeyDown={(e) => {
                               if (e.key === "Enter") {
                                   validateHandler()
                               }
                           }}
                           className="border-black border pl-1"/>
                </div>
                <button onClick={validateHandler}
                        className="bg-green-700 rounded p-2 text-white font-semibold m-auto flex pl-6 pr-6 mt-4">Validate
                </button>
            </div>
        </div>
    )
}

AddTodoPopUp.propTypes = {
    toggleDisplay: PropTypes.func.isRequired,
}