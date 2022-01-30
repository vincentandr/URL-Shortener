import React, { useState, useEffect, createContext } from "react";
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
  const [drawerState, setDrawer] = useState(false);

  useEffect(() => {
    dispatch(fetchProducts());
    dispatch(fetchCart());
  }, []);

  return (
    <div className="App">
      <Navbar totalItems={cart.length} onClickDrawer={setDrawer} />
      <CartContext.Provider value={{ cart: cart, onClickDrawer: setDrawer }}>
        <Cart drawerState={drawerState} />
        <Products products={products} />
      </CartContext.Provider>
    </div>
  );
}

export { App, CartContext };
