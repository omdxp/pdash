import { Component, createEffect, createSignal, onMount } from "solid-js";
import { Customer, Order } from "../../interfaces";

import DeleteModal from "../delete-modal";

interface ModifyCustomerModalProps {
  customer: Customer;
  submit: (draftCustomer: Customer) => void;
  cancel: () => void;
  deleteCustomer: () => void;
}

const ModifyCustomerModal: Component<ModifyCustomerModalProps> = ({
  customer,
  submit,
  cancel,
  deleteCustomer,
}) => {
  const [deleteModalShown, setDeleteModalShown] = createSignal(false);
  const [draftCustomer, setDraftCustomer] = createSignal(customer);
  onMount(() => {
    setDraftCustomer(customer);
  });
  return (
    <>
      <div
        id="authentication-modal"
        tabindex="-1"
        aria-hidden="true"
        class="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 w-full md:inset-0 h-modal md:h-full"
      >
        <div class="relative p-4 w-full max-w-md h-full md:h-auto">
          <div class="relative bg-white rounded-lg shadow dark:bg-gray-700">
            <button
              type="button"
              class="absolute top-3 right-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center dark:hover:bg-gray-800 dark:hover:text-white"
              onClick={cancel}
            >
              <svg
                aria-hidden="true"
                class="w-5 h-5"
                fill="currentColor"
                viewBox="0 0 20 20"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  fill-rule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clip-rule="evenodd"
                ></path>
              </svg>
              <span class="sr-only">Close modal</span>
            </button>
            <div class="py-6 px-6 lg:px-8">
              <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">
                Modify Customer
              </h3>
              <form class="space-y-6" action="#">
                <div>
                  <label
                    for="email"
                    class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300"
                  >
                    Name
                  </label>
                  <input
                    type="text"
                    class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                    placeholder="Name"
                    value={draftCustomer().name}
                    onChange={(e) =>
                      setDraftCustomer({
                        ...draftCustomer(),
                        name: (e.target as HTMLInputElement).value,
                      })
                    }
                  />
                </div>
                <button
                  type="submit"
                  onClick={() => submit(draftCustomer())}
                  class="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
                >
                  Modify
                </button>
              </form>
              <button
                onClick={() => setDeleteModalShown(true)}
                class="my-2 w-full text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
      {deleteModalShown() && (
        <DeleteModal
          submit={deleteCustomer}
          cancel={() => setDeleteModalShown(false)}
        />
      )}
    </>
  );
};
export default ModifyCustomerModal;
