import React, { useEffect, useState } from "react";
import { Todo } from "../components/Todo";
import { TodoForm } from "../components/TodoForm";
import { EditTodoForm } from "../components/EditTodoForm";
import { connect, useDispatch } from "react-redux";
import * as types from '../utils/constant/ActionTypes';
import api from "../utils/api";

const TodoWrapper = ({ todoState }) => {
    const [loading, setLoading] = useState(false);
    const dispatch = useDispatch();

    useEffect(() => {
        const fetchTodos = async () => {
            setLoading(true);
            try {
                const response = await api.get(`/get-tasks`);
                dispatch({
                    type: types.SET_TODOS,
                    payload: response?.data,
                });
            } catch (error) {
                console.error("Error fetching todos:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchTodos();
    }, [dispatch]);

    const addTodo = async (todo) => {
        setLoading(true);
        try {
            const response = await api.post(`/create-task`, {
                name: todo,
            });
            dispatch({
                type: types.ADD_TODO,
                payload: response?.data,
            });
        } catch (error) {
            console.error("Error while adding todos:", error);
        } finally {
            setLoading(false);
        }
    }

    const deleteTodo = async (id) => {
        setLoading(true);
        try {
            await api.delete(`/delete-task/${id}`);
            dispatch({
                type: types.DELETE_TODO,
                payload: {
                    id,
                }
            });
        } catch (error) {
            console.error("Error while deleting todo:", error);
        } finally {
            setLoading(false);
        }
    }

    const toggleComplete = async (id, completed) => {
        setLoading(true);
        try {
            await api.put(`/update-task-completed/${id}`, {
                completed: !completed,
            });
            dispatch({
                type: types.TASK_COMPLETED,
                payload: {
                    completed: !completed,
                    id,
                }
            });
        } catch (error) {
            console.error("Error while editing to completed todo:", error);
        } finally {
            setLoading(false);
        }
    }

    const editTodo = (id) => {
        dispatch({
            type: types.SET_EDIT_TODO,
            payload: {
                id,
            }
        });
    }

    const editTask = async (task, id) => {
        setLoading(true);
        try {
            await api.put(`/update-task/${id}`, {
                name: task,
            });
            dispatch({
                type: types.EDIT_TODO,
                payload: {
                    name: task,
                    id,
                }
            });
        } catch (error) {
            console.error("Error while editing todo:", error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="todo-container">
            <div className="TodoWrapper">
                <h1 className="todo-header">Get Things Done !</h1>
                {loading ? (
                    <div className="loader" />
                ) : (
                    <>
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
                    </>
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
