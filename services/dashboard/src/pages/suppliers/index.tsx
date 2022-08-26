import { Component, For } from "solid-js";

import { suppliers } from "../../store";

const Suppliers: Component = ({}) => {
  return (
    <>
      {suppliers.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          <For each={suppliers()}>
            {(supplier) => (
              <div>
                <h2>ID: {supplier.id}</h2>
                <p>Name: {supplier.name}</p>
                <p>
                  Created at: {new Date(supplier.created_at).toDateString()} -{" "}
                  {new Date(supplier.created_at).toTimeString()}
                </p>
                <p>
                  Updated at: {new Date(supplier.updated_at).toDateString()} -{" "}
                  {new Date(supplier.updated_at).toTimeString()}
                </p>
              </div>
            )}
          </For>
        </div>
      )}
    </>
  );
};
export default Suppliers;
