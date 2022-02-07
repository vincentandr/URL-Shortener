import { combineReducers } from "redux";

import { cart } from "./Cart";
import { products } from "./Products";
import { payment } from "./Payment";

const rootReducer = combineReducers({
  cart,
  products,
  payment,
});

export default rootReducer;
