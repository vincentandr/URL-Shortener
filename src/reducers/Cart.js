export const cart = (state = [], action) => {
  switch (action.type) {
    case "FETCH_CART":
      return action.payload;
    case "ADD_ITEM":
      return action.payload;
    case "REMOVE_ITEM":
      return action.payload;
    case "REMOVE_ALL":
      return action.payload;
    case "CHECKOUT":
    default:
      return state;
  }
};
