import React, { useState, useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Routes, Route, BrowserRouter as Router } from "react-router-dom";

import { Products, Navbar, Payment } from "../components";
import { fetchProducts, fetchCart } from "../actions";

const getSelectors = (state) => ({
  cart: state.cart,
  products: state.products,
});

function App() {
  const dispatch = useDispatch();

  const { cart, products } = useSelector(getSelectors);
  const [drawerState, setDrawer] = useState(false);
  const [loginState, setLogin] = useState(false);

  useEffect(() => {
    dispatch(fetchProducts());
    dispatch(fetchCart());
  }, []);

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route
            exact
            path="/"
            element={
              <>
                <Navbar
                  cart={cart}
                  drawer={{ state: drawerState, onClick: setDrawer }}
                  login={{ state: loginState, onClick: setLogin }}
                />
                <Products
                  products={products}
                  cart={cart}
                  onClickDrawer={setDrawer}
                />
              </>
            }
          />
          <Route
            exact
            path="/payment"
            element={<Payment cart={cart} />}
          ></Route>
        </Routes>
      </div>
    </Router>
  );
}

export { App };
