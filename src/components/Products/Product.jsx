import {React, useContext} from "react";
import { useDispatch } from "react-redux";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton} from "@mui/material"
import { AddShoppingCart } from "@mui/icons-material";

import { CartContext } from "../../pages/App"
import { addCartItem } from "../../actions";

const Product = ({product}) => {
    const dispatch = useDispatch()
    const cart = useContext(CartContext)

    const handleCartClick = (productId, qty) => {
        let obj = cart.find(item => item.product_id === productId)

        if (obj !== undefined){
            qty += obj.qty
        }

        dispatch(addCartItem(productId, qty))
    }

    return (
        <Card>
            <CardMedia image={product.image} title={product.name} style={{height:0, paddingTop: "56.25%"}}/>
            <CardContent>
                <div>
                    <Typography variant="h5" gutterBottom>
                        {product.name}
                    </Typography>
                    <Typography variant="h5">
                        {product.price}
                    </Typography>
                    <Typography variant="h5">
                        {product.qty}
                    </Typography>
                </div>
                <Typography variant="body2">{product.desc}</Typography>
            </CardContent>
            <CardActions>
                <IconButton aria-label="Add to cart" onClick={() => handleCartClick(product.product_id, 1)}>
                    <AddShoppingCart />
                </IconButton>
            </CardActions>
        </Card>
    )
}

export default Product;