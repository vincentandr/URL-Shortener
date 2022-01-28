import axios from "axios";
import config from "./config";

const fetchCart = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1";
    const { data } = await axios.get(uri);

    dispatch({ type: "FETCH_CART", payload: data.products });
  } catch (error) {
    console.log(error);
  }
};

const addCartItem = (productId, quantity) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/cart/user1/" + productId;

    const { data } = await axios.put(
      uri,
      {},
      {
        params: {
          qty: quantity,
        },
      }
    );

    dispatch({ type: "ADD_ITEM", payload: data.products });
  } catch (error) {
    console.log(error);
  }
};

export { fetchCart, addCartItem };
