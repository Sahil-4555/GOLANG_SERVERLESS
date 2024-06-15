import * as types from "../../utils/constant/ActionTypes"

const initalState = {
    Todo: [],
}

const TodoReducer = (state = initalState.Todo, action) => {
    switch (action.type) {
        case types.ADD_TODO:
            return [
                ...state,
                action.payload,
            ];

        case types.SET_TODOS:
            return [...action.payload];

        case types.SET_EDIT_TODO:
            return state.map((todo) =>
                todo.id === action.payload.id ? { ...todo, is_editing: !todo.is_editing } : todo
            )

        case types.TASK_COMPLETED:
            return state.map((todo) =>
                todo.id === action.payload.id ? { ...todo, completed: !todo.completed } : todo
            )

        case types.EDIT_TODO:
            return state.map((todo) => todo.id === action.payload.id ? { ...todo, is_editing: !todo.is_editing, name: action.payload.name } : todo)

        case types.DELETE_TODO:
            return state.filter((todo) => todo.id !== action.payload.id)

        case types.SET_LOADING:
            return {
                ...state,
                loading: true,
            };
            
        case types.UNSET_LOADING:
            return {
                ...state,
                loading: false,
            };

        default:
            return state;
    }
}

export default TodoReducer