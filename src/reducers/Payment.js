export const payment = (state = { order: {}, secret_key: "" }, action) => {
  switch (action.type) {
    case "FETCH_DRAFT_ORDER": {
      const data = action.payload;
      return {
        ...state,
        order: data.order,
        secret_key: data.client_secret,
      };
    }
    case "MAKE_PAYMENT": {
      return state;
    }
    default:
      return state;
  }
};
