import axios from "axios";

const fetchDraftOrder = (userId) => async (dispatch) => {
  try {
    const uri =
      process.env.REACT_APP_SERVER_API_URL + "/payment/draft/" + userId;
    const { data } = await axios.get(uri);

    dispatch({ type: "FETCH_DRAFT_ORDER", payload: data });
  } catch (error) {
    console.log(error);
  }
};

const makePayment = (order) => async (dispatch) => {
  try {
    const uri =
      process.env.REACT_APP_SERVER_API_URL + "/payment/" + order.order_id;
    const { data } = await axios.post(uri, {
      customer: order.customer,
    });

    dispatch({ type: "MAKE_PAYMENT", payload: data });
  } catch (error) {
    console.log(error);
  }
};

export { fetchDraftOrder, makePayment };
