import { Component, For } from "solid-js";

import { customers } from "../../store";

const Customers: Component = ({}) => {
  return (
    <>
      {customers.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          <For each={customers()}>
            {(customer) => (
              <div>
                <h2>ID: {customer.id}</h2>
                <p>Name: {customer.name}</p>
                <p>
                  Created at: {new Date(customer.created_at).toDateString()} -{" "}
                  {new Date(customer.created_at).toTimeString()}
                </p>
                <p>
                  Updated at: {new Date(customer.updated_at).toDateString()} -{" "}
                  {new Date(customer.updated_at).toTimeString()}
                </p>
              </div>
            )}
          </For>
        </div>
      )}
    </>
  );
};
export default Customers;
