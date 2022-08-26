import { Component, For, createEffect } from "solid-js";

import { orders } from "../../store";

const Home: Component = ({}) => {
  return (
    <>
      {orders.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          <For each={orders()}>
            {(order) => (
              <div>
                <h2>ID: {order.id}</h2>
                <p>Total price: {order.total_price}</p>
                <p>Customer ID: {order.customer_id}</p>
                <p>Supplier ID: {order.supplier_id}</p>
                <p>
                  Created at: {new Date(order.created_at).toDateString()} -{" "}
                  {new Date(order.created_at).toTimeString()}
                </p>
                <p>
                  Updated at: {new Date(order.updated_at).toDateString()} -{" "}
                  {new Date(order.updated_at).toTimeString()}
                </p>
              </div>
            )}
          </For>
        </div>
      )}
    </>
  );
};
export default Home;
