<?php
function post_to_go($type, $action, $data, $headers) {
   $fields = '';

   foreach($data as $key => $value) { 
      $fields .= $key . '=' . $value . '&'; 
   }

   rtrim($fields, '&');

   $post = curl_init();

   curl_setopt($post, CURLOPT_URL, GO_PATH.$type.'/'.$action );
   curl_setopt($post, CURLOPT_POST, count($data));
   curl_setopt($post, CURLOPT_POSTFIELDS, $fields);
   curl_setopt($post, CURLOPT_RETURNTRANSFER, 1);
   curl_setopt($post, CURLOPT_HTTPHEADER, $headers);

   $result = curl_exec($post);

   curl_close($post);

   return $result;
}