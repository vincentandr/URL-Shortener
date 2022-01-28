export const products = (state = [], action) => {
  switch (action.type) {
    case "FETCH_PRODUCTS": {
      return action.payload;
    }
    case "SEARCH_PRODUCTS": {
      return action.payload;
    }
    default:
      return state;
  }
};
