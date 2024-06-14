import React, { useDebugValue, useEffect, useState } from "react";
import { Todo } from "../components/Todo";
import { TodoForm } from "../components/TodoForm";
import { v4 as uuidv4 } from "uuid";
import { EditTodoForm } from "../components/EditTodoForm";
import axios from "axios";
import { connect, useDispatch } from "react-redux";
import * as types from '../utils/constant/ActionTypes'

const TodoWrapper = ({ todoState }) => {
    const token = localStorage.getItem("token")
    const dispatch = useDispatch();

    useEffect(() => {
        const fetchTodos = async () => {
            try {
                const response = await axios.get(`${process.env.URL}/get-tasks`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });

                dispatch({
                    type: types.SET_TODOS,
                    payload: response?.data,
                })


            } catch (error) {
                console.error("Error fetching todos:", error);
            }
        };

        fetchTodos();
    }, [])

    const addTodo = async (todo) => {
        try {
            const response = await axios.post(`${process.env.URL}/create-task`, {
                name: todo,
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            dispatch({
                type: types.ADD_TODO,
                payload: response?.data,
            })


        } catch (error) {
            console.error("Error while adding todos:", error);
        }
    }

    const deleteTodo = async (id) => {
        try {
            const response = await axios.delete(`${process.env.URL}/delete-task/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            dispatch({
                type: types.DELETE_TODO,
                payload: {
                    id,
                }
            })
        } catch (error) {
            console.error("Error while deleting todo:", error);
        }
    }

    const toggleComplete = async (id, completed) => {
        try {
            const response = await axios.put(`${process.env.URL}/update-task-completed/${id}`, {
                completed: !completed,
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            dispatch({
                type: types.TASK_COMPLETED,
                payload: {
                    completed: completed,
                    id,
                }
            })
        } catch (error) {
            console.error("Error while editing to completed todo:", error);
        }
    }

    const editTodo = (id) => {
        dispatch({
            type: types.SET_EDIT_TODO,
            payload: {
                id,
            }
        })
    }

    const editTask = async (task, id) => {
        try {
            const response = await axios.put(`${process.env.URL}/update-task/${id}`, {
                name: task,
            }, {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            dispatch({
                type: types.EDIT_TODO,
                payload: {
                    name: task,
                    id,
                }
            })
        } catch (error) {
            console.error("Error while editing todo:", error);
        }
    };

    return (
        <div className="todo-container">
            <div className="TodoWrapper">
                <h1 className="todo-header">Get Things Done !</h1>
                <TodoForm addTodo={addTodo} />
                {todoState.map((todo) =>
                    todo.is_editing ? (
                        <EditTodoForm editTodo={editTask} task={todo} />
                    ) : (
                        <Todo
                            key={todo.id}
                            task={todo}
                            deleteTodo={deleteTodo}
                            editTodo={editTodo}
                            toggleComplete={toggleComplete}
                        />
                    )
                )}
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    todoState: state?.TodoReducer || [],
});

export default connect(
    mapStateToProps
)(TodoWrapper);
