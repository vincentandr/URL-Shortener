import axios from "axios";
import config from "./config";

const fetchProducts = () => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/products";
    const { data } = await axios.get(uri);

    dispatch({ type: "FETCH_PRODUCTS", payload: data.products });
  } catch (error) {
    console.log(error);
  }
};

const searchProducts = (input) => async (dispatch) => {
  try {
    const uri = config.apiUrl + "/products/search";
    const { data } = await axios.get(uri, {
      params: {
        name: input,
      },
    });

    let result = [];

    if (Object.keys(data).length !== 0) {
      result = data.products;
    }

    dispatch({ type: "FETCH_PRODUCTS", payload: result });
  } catch (error) {
    console.log(error);
  }
};

export { fetchProducts, searchProducts };
