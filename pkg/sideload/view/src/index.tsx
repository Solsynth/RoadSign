import "./index.css";

/* @refresh reload */
import { render } from "solid-js/web";

import { Route, Router } from "@solidjs/router";

import RootLayout from "./layouts/RootLayout";
import Dashboard from "./pages/dashboard";

const root = document.getElementById("root");

render(() => (
    <Router root={RootLayout}>
        <Route path="/" component={Dashboard} />
    </Router>
), root!);
