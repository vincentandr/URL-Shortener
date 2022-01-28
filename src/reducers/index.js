import { combineReducers } from "redux";

import { cart } from "./Cart";
import { products } from "./Products";

const rootReducer = combineReducers({
  cart,
  products,
});

export default rootReducer;
