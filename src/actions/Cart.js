import axios from "axios";
import config from "./config";

const fetchCart = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1";
    const { data } = await axios.get(uri);

    let result = [];

    if (Object.keys(data).length !== 0) {
      result = data.products;
    }

    dispatch({ type: "FETCH_CART", payload: result });
  } catch (error) {
    console.log(error);
  }
};

const addCartItem = (item) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1/" + item.productId;

    const { data } = await axios.put(
      uri,
      {},
      {
        params: {
          name: item.name,
          qty: item.qty,
          price: item.price,
          desc: item.desc,
          image: item.image,
        },
      }
    );

    dispatch({ type: "ADD_ITEM", payload: data.products });
  } catch (error) {
    console.log(error);
  }
};

const removeCartItem = (productId) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1/" + productId;

    const { data } = await axios.delete(uri);

    let result = [];

    if (Object.keys(data).length !== 0) {
      result = data.products;
    }

    dispatch({ type: "REMOVE_ITEM", payload: result });
  } catch (error) {
    console.log(error);
  }
};

const removeAllCartItems = (item) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1";

    const { data } = await axios.delete(uri);

    let result = [];

    if (Object.keys(data).length !== 0) {
      result = data.products;
    }

    dispatch({ type: "REMOVE_ALL", payload: result });
  } catch (error) {
    console.log(error);
  }
};

export { fetchCart, addCartItem, removeCartItem, removeAllCartItems };
