export const cart = (state = { products: [], subtotal: 0 }, action) => {
  switch (action.type) {
    case "FETCH_CART": {
      const data = action.payload;
      let newProducts = [...state.products];
      newProducts = data.products;
      return {
        ...state,
        products: newProducts,
        subtotal: data.subtotal,
      };
    }
    case "ADD_ITEM": {
      const data = action.payload;
      let newProducts = [...state.products];
      newProducts = data.products;
      return {
        ...state,
        products: newProducts,
        subtotal: data.subtotal,
      };
    }
    case "REMOVE_ITEM": {
      const data = action.payload;
      let newProducts = [...state.products];
      newProducts = data.products;
      return {
        ...state,
        products: newProducts,
        subtotal: data.subtotal,
      };
    }
    case "REMOVE_ALL": {
      const data = action.payload;
      let newProducts = [...state.products];
      newProducts = data.products;
      return {
        ...state,
        products: newProducts,
        subtotal: data.subtotal,
      };
    }
    case "CHECKOUT": {
    }
    default:
      return state;
  }
};
