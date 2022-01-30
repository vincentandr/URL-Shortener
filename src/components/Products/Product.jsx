import React from "react"
import {useContext} from "react";
import { useDispatch } from "react-redux";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton} from "@mui/material"
import { AddShoppingCart } from "@mui/icons-material";

import { CartContext } from "../../pages/App"
import { addCartItem } from "../../actions";

const Product = ({product}) => {
    const dispatch = useDispatch()
    const value = useContext(CartContext)

    const handleCartClick = (productId, qty, price, desc, image) => {
        // Add +1 to qty if item already exists in cart
        let obj = value.cart.find(item => item.product_id === productId)

        if (obj !== undefined){
           qty += obj.qty
        }

        let item = {
            productId: productId,
            qty: qty,
            price: price,
            desc: desc,
            image: image,
        }

        dispatch(addCartItem(item)).then((result) => {
            value.onClickDrawer(true)
        })
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
                <IconButton aria-label="Add to cart" onClick={() => handleCartClick(
                    product.product_id,
                    1,
                    product.price,
                    product.desc,
                    product.image)}>
                    <AddShoppingCart />
                </IconButton>
            </CardActions>
        </Card>
    )
}

export default React.memo(Product);