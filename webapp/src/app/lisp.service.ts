import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class LispService {

  examples = {
    addition: "(+ 1 2)",
    comparison: "(> 1 0)",
    condition: "(if (> 1 0) 1 0)",
    createList: "(setq toto (+ 1 2))",
    immutableList: "(cons 1 (cons 2 (cons 3 NIL)))",
    functionDefinitionWithComments: "(defun multiply_by_seven (number) \"Multiply NUMBER by seven.\" (* 7 number))",
    highOrderFunction: "(defun multiply_by_seven (number) (* 7 number)(* 7 number))\n(multiply_by_seven 8)",
    simpleRecursive: "(defun add (a b)\n" +
      "  (if (> a 0)\n" +
      "      (add (- a 1) (+ b 1))\n" +
      "      b\n" +
      "  )\n" +
      ")\n" +
      "(add 10 5)",
    recursiveFunction: "(defun list_length_ (list n)\n" +
      "  (if list (list_length_ (cdr list) (+ n 1) ) n)\n" +
      ")\n" +
      "(defun list_length (list) \n" +
      "  (list_length_ list 0)\n" +
      ")\n" +
      "(list_length '(1 2 3 4))"}

  constructor(private http: HttpClient) { }

  evalLispCode(lispCode: string) {

    const formData = new FormData();
    formData.append('expression', lispCode);

    return this.http.post("http://localhost:8080/lispgo", formData, {})
  }

  getExamples = () => {
    return this.examples;
  }

}
