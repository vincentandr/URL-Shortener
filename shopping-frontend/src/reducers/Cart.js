export const cart = (
  state = { products: [], subtotal: 0, order_id: "" },
  action
) => {
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
      const data = action.payload;
      return {
        ...state,
        order_id: data.order_id,
      };
    }
    default:
      return state;
  }
};
