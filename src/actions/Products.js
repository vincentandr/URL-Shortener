import axios from "axios";

const fetchProducts = () => async (dispatch) => {
  try {
    const uri = process.env.REACT_APP_SERVER_API_URL + "/products";
    const { data } = await axios.get(uri);

    dispatch({ type: "FETCH_PRODUCTS", payload: data.products });
  } catch (error) {
    console.log(error);
  }
};

const searchProducts = (input) => async (dispatch) => {
  try {
    const uri = process.env.REACT_APP_SERVER_API_URL + "/products/search";
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
