import { Route, Routes } from "@solidjs/router";

import type { Component } from "solid-js";
import Customer from "./pages/customer";
import Customers from "./pages/customers";
import Header from "./components/header";
import Home from "./pages/home";
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
      </Routes>
    </>
  );
};

export default App;