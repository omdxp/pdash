import { Component, createEffect, onMount } from "solid-js";
import { Route, Routes, useNavigate } from "@solidjs/router";
import { setToken, token } from "./store";

import Customer from "./pages/customer";
import Customers from "./pages/customers";
import Header from "./components/header";
import Home from "./pages/home";
import Login from "./pages/login";
import Signup from "./pages/signup";
import Supplier from "./pages/supplier";
import Suppliers from "./pages/suppliers";

const App: Component = () => {
  return (
    <>
      <Header />
      <Routes>
        <Route path="/" component={Home} />
        <Route path="/suppliers" component={Suppliers} />
        <Route path="/suppliers/:id" component={Supplier} />
        <Route path="/customers" component={Customers} />
        <Route path="/customers/:id" component={Customer} />
        <Route path="/login" component={Login} />
        <Route path="/signup" component={Signup} />
      </Routes>
    </>
  );
};

export default App;
