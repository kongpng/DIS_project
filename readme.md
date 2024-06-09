Introduction

We chose to make a delivery service/food ordering application that would track different prices and compare them. We took a bit of a detour and ended up creating more of a back-end tool, that shows the history of orders created between different delivery services and restaurants. This was because we chose to add customers to our database and create 'orders' that connect the customer with the delivery service and the restaurant. We chose to focus on four different entities, namely customer, deliveryservice, restaurant and dishes. Each of said entities would have varying attributes, the chief among them being an ID serving as our primary key to easily identify each entity internally. The 'order' relation serves as a history with the ID of a single customer, delivery service and restaurant, as well as a date of said order. This allows a user of the database to easily search up where a customer placed an order at a given time. Each order only has a single customer, delivery and restaurant. The 'menu' relation between a restaurants and its dishes contains only a single restaurant, but several dishes' IDs. We chose to specify the type for each attribute to ensure clarity. With for example the 'e-mail' attribute of the customer, we chose to not make said attribute a key, because we could see a situation where a user deletes their account and creates a new one with the same email. In that situation we reckoned that they should be considered a new user, with a new ID.


We explored the possibility of implementing regular expression matching and context-free as required by the assignment specifications. Ultimately, we failed at this, due to us starting extraordinarily late due to some personal issues among all of us, and switching to SQLite at the very last minute. We sadly ended up discovering (too late) that SQLite did not have support for the regex function natively like postgres does, so what we really just ended up with is a rather simple CRUD application. We hope that this doesn't end up with us failing the subject, but we also realize we should've started way earlier. 

How to run the code / README

When in the directory you should intialise the database by running the following lines in your terminal. (The following is case sensitive)
\begin{verbatim}
    cd db
    sqlite3 Data.db
    .read init.sql
    .read initial_data.sql
    ** PRESS CONTROL C TO EXIT SQLITE**
    cd ..
    go run project
\end{verbatim}
This will host the webapp at http://localhost:8080/ , which can be accessed from your browser.

\textbf{Prerequisites}
\begin{itemize}
    \item GCC
    \item Sqlite
    \item Go
\end{itemize}
