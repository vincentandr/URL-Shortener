import axios from "axios";
import config from "./config";

const fetchCart = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1";
    const { data } = await axios.get(uri);

    if (Object.keys(data).length === 0) {
      data["products"] = [];
      data["subtotal"] = 0;
    }

    dispatch({ type: "FETCH_CART", payload: data });
  } catch (error) {
    console.log(error);
  }
};

const addCartItem = (productId, qty) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1/" + productId;

    const { data } = await axios.put(
      uri,
      {},
      {
        params: {
          qty: qty,
        },
      }
    );

    dispatch({ type: "ADD_ITEM", payload: data });
  } catch (error) {
    console.log(error);
  }
};

const removeCartItem = (productId) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1/" + productId;

    const { data } = await axios.delete(uri);

    if (Object.keys(data).length === 0) {
      data["products"] = [];
      data["subtotal"] = 0;
    }

    dispatch({ type: "REMOVE_ITEM", payload: data });
  } catch (error) {
    console.log(error);
  }
};

const removeAllCartItems = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1";

    const { data } = await axios.delete(uri);

    if (Object.keys(data).length === 0) {
      data["products"] = [];
      data["subtotal"] = 0;
    }

    dispatch({ type: "REMOVE_ALL", payload: data });
  } catch (error) {
    console.log(error);
  }
};

const checkout = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/checkout/user1";

    const { data } = await axios.get(uri);

    dispatch({ type: "CHECKOUT", payload: data });
  } catch (error) {
    console.log(error);
  }
};

export { fetchCart, addCartItem, removeCartItem, removeAllCartItems, checkout };
