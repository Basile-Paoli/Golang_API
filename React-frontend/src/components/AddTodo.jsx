import {memo, useEffect, useRef, useState} from "react";
import {AddTodoPopUp} from "./AddTodoPopUp.jsx";

export function AddTodo() {
    const popUpRef = useRef(null);
    const [displayPopUp, setDisplayPopUp] = useState(false);
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (popUpRef.current && !popUpRef.current.contains(event.target)) {
                setDisplayPopUp(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return () => {
            document.removeEventListener('mousedown', handleClickOutside)
        }
    }, []);
    return (
        <>
            <button autoFocus onClick={() => setDisplayPopUp(true)}
                    className="bg-green-700 rounded p-2 text-white font-semibold m-auto flex pl-6 pr-6">
                Add todo
            </button>
            {displayPopUp && <AddTodoPopUp toggleDisplay={setDisplayPopUp}/>}
        </>
    )
}

const AddTodoMemo = memo(AddTodo);
export default AddTodoMemo


