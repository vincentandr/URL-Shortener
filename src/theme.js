import { Button, Typography } from "@mui/material";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import { Link } from "react-router-dom";

const blackButtonTheme = createTheme({
  palette: {
    primary: {
      main: "#212121",
    },
    action: {
      disabledBackground: "black",
      disabled: "inherit",
    },
  },
});

const logoFontTheme = createTheme({
  palette: {
    primary: {
      main: "#212121",
    },
  },
  typography: {
    fontFamily: "Fugaz One",
  },
});

const BlackButton = ({
  component,
  to,
  disabled,
  type,
  variant,
  text,
  onClick,
  fullWidth,
}) => (
  <ThemeProvider theme={blackButtonTheme}>
    <Button
      component={component}
      to={to}
      disabled={disabled}
      type={type}
      variant={variant}
      onClick={onClick}
      fullWidth={fullWidth}
    >
      {text}
    </Button>
  </ThemeProvider>
);

const Logo = ({ variant }) => (
  <ThemeProvider theme={logoFontTheme}>
    <Typography
      component={Link}
      to="/"
      color="inherit"
      variant={variant}
      sx={{
        textDecoration: "none",
      }}
    >
      Microshopping
    </Typography>
  </ThemeProvider>
);

export { BlackButton, Logo };
