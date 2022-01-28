import React, { useEffect, createContext } from "react";
import { useSelector, useDispatch } from "react-redux";
import "../css/App.css";

import { Products, Navbar, Cart } from "../components";
import { fetchProducts, fetchCart } from "../actions";

const getSelectors = (state) => ({
  cart: state.cart,
  products: state.products,
});

const CartContext = createContext();

function App() {
  const dispatch = useDispatch();

  const { cart, products } = useSelector(getSelectors);

  useEffect(() => {
    dispatch(fetchProducts());
    dispatch(fetchCart());
  }, []);

  return (
    <div className="App">
      <Navbar totalItems={cart.length} />
      {/* <Cart cart={cart}></Cart> */}
      <CartContext.Provider value={cart}>
        <Products products={products} />
      </CartContext.Provider>
    </div>
  );
}

export { App, CartContext };
